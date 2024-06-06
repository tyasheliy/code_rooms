<template>
  <el-row class="main-row" align="middle">
    <el-col :span="8"></el-col>
    <el-col :span="8">
      <el-card>
        <template #header>
          <div class="card-header">
            <span>Регистрация</span>
          </div>
        </template>
        <el-form>
          <el-form-item>
            <el-alert
              :closable="false"
              :title="alertMessage"
              :type="alertType"
              center
              show-icon
            />
          </el-form-item>
          <el-form-item>
            <el-input v-model="login" placeholder="Логин:" />
          </el-form-item>
          <el-form-item>
            <el-input v-model="password" show-password placeholder="Пароль:" />
          </el-form-item>
          <el-form-item>
            <el-input
              v-model="confirmPassword"
              show-password
              placeholder="Повторите пароль:"
            />
          </el-form-item>
          <el-form-item>
            <div class="submit_buttons">
              <el-button @click="onSubmit">Зарегистрироваться</el-button>
              <router-link to="/login" class="el-button">Назад</router-link>
            </div>
          </el-form-item>
        </el-form>
      </el-card>
    </el-col>
    <el-col :span="8"></el-col>
  </el-row>
</template>

<script setup>
import { ref } from "vue";
import axios from "axios";

const login = ref("");
const password = ref("");
const confirmPassword = ref("");

const alertMessage = ref("");
const alertType = ref("");

const onSubmit = () => {
  if (login.value === "" || password.value === "") {
    showError("Заполните поля!");
    return;
  }

  const data = {
    login: login.value,
    password: password.value,
  };

  axios
    .post("/api/v1/auth/register", data)
    .then(() => {
      showSuccess("Вы успешно зарегистрированы!");
    })
    .catch((err) => {
      if (err.response.data.message) {
        showError(err.response.data.message);
      } else {
        showError(err);
      }
    });
};

const showSuccess = (message) => {
  console.log(message);
  alertMessage.value = message;
  alertType.value = "success";
};

const showError = (message) => {
  console.error(message);
  alertMessage.value = message;
  alertType.value = "error";
};
</script>

<style scoped>
.submit_buttons {
  height: 100%;
  width: 100%;
  display: flex;
  justify-content: center;
  flex-direction: column;
}

.submit_buttons .el-button {
  margin: 3px;
}

a {
  text-decoration: none;
}

.main-row {
  height: 70vh;
  width: 100%;
}
</style>
