<template>
  <router-view v-slot="{ Component }">
    <keep-alive>
      <component :is="Component" />
    </keep-alive>
  </router-view>
</template>
<script lang="ts">
import { defineComponent, onMounted } from "vue";
import { useStore } from "vuex";

export default defineComponent({
  name: "App",
  setup() {
    const $store = useStore();

    onMounted(() => {
      $store.dispatch("getConfig");
    });

    return {};
  }
});
</script>

<style lang="stylus"></style>
