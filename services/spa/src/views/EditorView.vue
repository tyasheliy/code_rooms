<template>
  <el-dialog
    title="Выберите действие"
    v-model="actionDialogVisible"
    :show-close="false"
  >
    <el-row class="action-dialog-row">
      <el-col :span="12" class="text-center">
        <el-button
          @click="
            entry = 'new';
            actionDialogVisible = false;
            handleConnection();
          "
        >
          Создать новую сессию
        </el-button>
      </el-col>
      <el-col :span="12" class="text-center">
        <el-button
          @click="
            actionDialogVisible = false;
            entryInputDialogVisible = true;
          "
        >
          Войти по коду приглашения
        </el-button>
      </el-col>
    </el-row>
  </el-dialog>

  <el-dialog
    v-model="entryInputDialogVisible"
    title="Вход по коду приглашения"
    :show-close="false"
  >
    <el-form>
      <el-form-item label="Код: ">
        <el-input v-model="entry" show-password />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button
        type="primary"
        @click="
          entryInputDialogVisible = false;
          handleConnection();
        "
        >Принять</el-button
      >
    </template>
  </el-dialog>

  <el-container class="full-size-container">
    <el-aside width="15%">
      <el-container>
        <el-main>
          <el-menu>
            <el-scrollbar>
              <el-menu-item
                v-for="file in files"
                :key="file.name"
                @click="changeFile(file.name)"
              >
                {{ file.name }}
              </el-menu-item>
            </el-scrollbar>
          </el-menu>
        </el-main>
        <el-footer>
          <el-col class="action-buttons-wrapper">
            <el-row :span="12">
              <el-button
                class="action-button"
                type="primary"
                @click="createFileDialogVisible = true"
                >Новый файл</el-button
              >
            </el-row>
            <el-row :span="12">
              <el-button
                class="action-button"
                type="primary"
                @click="inviteDialogVisible = true"
              >
                Пригласить
              </el-button>
            </el-row>
          </el-col>
          <el-dialog
            title="Пригласить участника"
            :align-center="true"
            v-model="inviteDialogVisible"
          >
            <span
              >Код для приглашения:
              <el-button link type="primary" @click="copyEntry">{{
                entry
              }}</el-button></span
            >
          </el-dialog>
          <el-dialog
            title="Ошибка"
            :align-center="true"
            v-model="errorDialogVisible"
          >
            <span>{{ errorDialogMessage }}</span>
          </el-dialog>
          <el-dialog
            title="Создание файла"
            :align-center="true"
            v-model="createFileDialogVisible"
          >
            <el-form>
              <el-form-item label="Имя файла: ">
                <el-input
                  v-model="fileDialogFilename"
                  class="file-dialog-input"
                />
              </el-form-item>
            </el-form>
            <template #footer>
              <div class="dialog-footer">
                <el-button @click="createFileDialogVisible = false"
                  >Отмена</el-button
                >
                <el-button type="primary" @click="createFile">
                  Принять
                </el-button>
              </div>
            </template>
          </el-dialog>
        </el-footer>
      </el-container>
    </el-aside>
    <el-main>
      <CodeEditor
        border-radius="0"
        :lang-list-display="true"
        :line-nums="true"
        :copy-code="false"
        :languages="languages"
        width="100%"
        height="100%"
        v-model="code"
      ></CodeEditor>
    </el-main>
  </el-container>
</template>

<style scoped>
.file-dialog-input {
  width: 40%;
}

.action-dialog-row {
  height: 40%;
}

.text-center {
  text-align: center;
}

.action-buttons-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.action-button {
  margin: 3px;
}
</style>

<script setup>
// eslint-disable-next-line no-unused-vars
import hljs from "highlight.js";
import CodeEditor from "simple-code-editor";
import { ref } from "vue";

import axios from "axios";
import { useRouter } from "vue-router";

const router = useRouter();

if (localStorage.getItem("access_token") === null) {
  router.push({ path: "/404" });
}

const actionDialogVisible = ref(true);
const entryInputDialogVisible = ref(false);
const entry = ref("");

const files = ref([]);

const languages = [
  ["go", "Go"],
  ["python", "Python"],
  ["php", "PHP"],
  ["html", "HTML"],
  ["javascript", "JS"],
];

let codeRecord = "";
const requestInterval = 1500;

const code = ref(codeRecord);

const inviteDialogVisible = ref(false);

const createFileDialogVisible = ref(false);
const fileDialogFilename = ref("");

const errorDialogVisible = ref(false);
const errorDialogMessage = ref("");

let isHandlingRequest = false;

const copyEntry = () => {
  navigator.clipboard.writeText(entry.value);
};

function showErrorDialog(message) {
  errorDialogMessage.value = message;
  errorDialogVisible.value = true;
}

// eslint-disable-next-line no-unused-vars
function createFile() {
  const filename = fileDialogFilename.value;

  createFileDialogVisible.value = false;

  if (filename === "") {
    showErrorDialog("Имя файла не может быть пустым");
    return;
  }

  if (files.value.find((f) => f.name === filename) !== undefined) {
    showErrorDialog("Файл с таким именем уже существует");
    return;
  }

  files.value.push({
    name: filename,
    data: "",
    current: false,
  });

  createFileDialogVisible.value = false;
}

function changeFile(filename) {
  let file = files.value.find((el) => el.name === filename);

  files.value.forEach((f) => {
    if (f.current) {
      f.current = false;
    }
  });

  file.current = true;
  handleFileChange(file);
}

function handleFileChange(file) {
  code.value = file.data;
}

function connect(e) {
  const connUrl = "ws://" + window.location.host + "/api/v1/editor/socket/" + e;
  const conn = new WebSocket(connUrl);

  conn.onopen = (ev) => {
    console.log("Connection successfully opened");
    console.log(ev);
  };

  conn.onclose = (ev) => {
    console.log(ev);

    showErrorDialog(
      "Не удалось подключится, либо подключение оборвано. Попробуйте подключиться снова."
    );
  };

  conn.onerror = (ev) => {
    console.log(ev);
  };

  conn.onmessage = (ev) => {
    const mess = JSON.parse(ev.data);

    let fileExists = false;

    files.value.forEach((f) => {
      if (mess.filename === f.name) {
        fileExists = true;
      }
    });

    if (fileExists) {
      const file = files.value.find((f) => f.name === mess.filename);

      file.data = mess.data;

      if (file.current) {
        handleFileChange(file);
      }
    } else {
      files.value.push({
        name: mess.filename,
        data: mess.data,
        current: false,
      });
    }

    console.log(files);
  };

  const handleCodeRequest = () => {
    if (isHandlingRequest) {
      return;
    }

    isHandlingRequest = true;

    if (codeRecord === code.value) {
      return;
    }

    let currentFilename;

    files.value.forEach((f) => {
      if (f.current) {
        currentFilename = f.name;
      }
    });

    const data = {
      filename: currentFilename,
      data: code.value,
    };

    console.log(files.value);
    console.log(data);
    conn.send(JSON.stringify(data));

    codeRecord = code.value;
  };

  document.addEventListener("keyup", () => {
    setTimeout(handleCodeRequest, requestInterval);
    isHandlingRequest = false;
  });
}

function handleConnection() {
  if (entry.value === "new") {
    console.log("creating new session");
    axios.post("/api/v1/sessions").then((r) => {
      console.log(r.data.id);
      axios.post("/api/v1/entries", { session_id: r.data.id }).then((r) => {
        console.log(r.data.id);
        entry.value = r.data.id;
        connect(r.data.id);
      });
    });
  } else {
    console.log("joining existing session");
    connect(entry.value);
  }
}
</script>
