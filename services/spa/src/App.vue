<template>
  <div class="main-container">
    <Navbar />
    <router-view />
  </div>
</template>

<script setup>
// eslint-disable-next-line no-unused-vars
import Navbar from "@/components/NavBar.vue";
</script>

<style>
html,
body {
  margin: 0;
  height: 100%;
  width: 100%;
}

.main-container {
  display: flex;
  flex-direction: column;
  position: relative;
  width: 100%;
  height: 100%;
}

#app {
  width: 100%;
  height: 100%;
}

body {
  font-family: "JetBrains Mono", "serif" !important;
}

button,
input,
select,
textarea {
  font-family: inherit;
  font-size: inherit;
  line-height: inherit;
  color: inherit;
}
</style>
