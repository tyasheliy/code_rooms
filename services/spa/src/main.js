import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import ElementPlus from "element-plus";
import "@fontsource/jetbrains-mono";
import "element-plus/dist/index.css";
import axios from "axios";

const app = createApp(App);

axios.defaults.baseURL = "http://" + window.location.host + "/";

app.use(ElementPlus);

app.use(router).mount("#app");
