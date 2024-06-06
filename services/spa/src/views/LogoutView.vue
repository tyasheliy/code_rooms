<script>
export default {
  mounted() {
    localStorage.clear();
    this.$router.push({ path: "/login" });
  },
};
</script>

<template>
  <div></div>
</template>
