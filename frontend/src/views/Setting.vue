<template>
  <header class="m-2 flex justify-end items-center">
    <router-link to="/" class="flex item-center text-blue-500 text-xl">
      <RollbackOutlined />
    </router-link>
  </header>
  <main class="m-2"></main>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { RollbackOutlined } from "@ant-design/icons-vue";
export default defineComponent({
  name: "Setting",
  components: { RollbackOutlined },
  setup() {
    return {};
  }
});
</script>

<style scoped></style>
