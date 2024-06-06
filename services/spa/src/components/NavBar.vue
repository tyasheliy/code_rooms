<template>
  <el-header>
    <el-menu mode="horizontal">
      <el-menu-item>
        <router-link to="/">Главная</router-link>
      </el-menu-item>
      <div class="flex-grow"></div>
      <el-menu-item>
        <router-link :to="{ path: logged ? '/logout' : '/login' }">
          {{ logged ? "Выйти" : "Войти" }}
        </router-link>
      </el-menu-item>
    </el-menu>
  </el-header>
</template>

<script>
export default {
  data() {
    return {
      logged: false,
    };
  },
  updated() {
    if (localStorage.getItem("access_token") !== null) {
      this.logged = true;
    }
  },
};
</script>

<style scoped></style>
