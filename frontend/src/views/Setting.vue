<template>
  <header class="p-3 flex justify-end items-center">
    <router-link to="/" class="flex item-center text-blue-500 text-xl">
      <RollbackOutlined />
    </router-link>
  </header>
  <main class="my-2 px-4">
    <div class="text-gray-700 pb-2">文件目录</div>
    <div
      @click="setOutDir"
      style="min-width: 200px;height: 32px;font-size: 15px;"
      class="inline-flex flex-auto bg-gray-100 rounded py-1 px-3 mb-2 text-gray-400 tracking-wider font-medium cursor-pointer hover:bg-gray-200 hover:text-gray-500"
    >
      {{ config.outDir }}
    </div>
    <a-form :model="config" layout="vertical" size="small">
      <a-form-item label="文件类型">
        <a-select
          v-model:value="config.target"
          style="width: 200px;"
          placeholder="选择文件类型"
        >
          <a-select-option value="jpg">jpg</a-select-option>
          <a-select-option value="png">png</a-select-option>
          <a-select-option value="webp">webp</a-select-option>
        </a-select>
      </a-form-item>
    </a-form>
  </main>
</template>

<script lang="ts">
import { computed, defineComponent, watch } from "vue";
import { RollbackOutlined } from "@ant-design/icons-vue";
import { useStore } from "vuex";
import { useDebounceFn } from "../composables/useDebounceFn";

export default defineComponent({
  name: "Setting",
  components: { RollbackOutlined },
  setup() {
    const $store = useStore();
    const config = computed<AppConfig>(() => {
      return $store.getters.config;
    });

    // 选择输出目录
    const setOutDir = async () => {
      const outDir = await window.backend.Config.SetOutDir();
      if (outDir) config.value.outDir = outDir;
    };

    // 保存配置
    const setConfig = useDebounceFn((v: AppConfig) => {
      $store.dispatch("setConfig", v);
    }, 500);

    // config有变动，则更新配置
    watch([() => config.value.outDir, () => config.value.target], (v, old) => {
      console.log(v, old);
      if (v && old && JSON.stringify(v) !== JSON.stringify(old)) {
        setConfig(config.value);
      }
    });

    return {
      config,
      setOutDir
    };
  }
});
</script>

<style scoped></style>
