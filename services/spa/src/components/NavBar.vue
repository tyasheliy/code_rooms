<template>
  <el-header>
    <el-menu mode="horizontal">
      <el-menu-item>
        <router-link to="/">Главная</router-link>
      </el-menu-item>
      <el-menu-item>
        <router-link to="/editor">Редактор</router-link>
      </el-menu-item>
      <el-menu-item>
        <router-link to="/login">Войти</router-link>
      </el-menu-item>
    </el-menu>
  </el-header>
</template>

<style scoped>
a {
  text-decoration: none;
}
</style>
