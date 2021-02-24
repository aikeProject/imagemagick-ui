<template>
  <header class="p-3 flex justify-end items-center">
    <router-link to="/" class="flex item-center text-blue-500 text-xl">
      <RollbackOutlined />
    </router-link>
  </header>
  <main class="my-2 px-4">
    <a-form :model="config" class="space-y-3" layout="vertical">
      <a-form-item label="文件目录" class="mb-0">
        <div
          @click="setOutDir"
          style="min-width: 200px;height: 32px;font-size: 15px;"
          class="inline-flex flex-auto bg-gray-100 rounded py-1 px-3 text-gray-400 tracking-wider font-medium cursor-pointer hover:bg-gray-200 hover:text-gray-500"
        >
          {{ config.outDir }}
        </div>
      </a-form-item>
      <a-form-item label="文件类型">
        <a-select
          class="item"
          v-model:value="config.target"
          placeholder="选择文件类型"
        >
          <a-select-option value="jpg">jpg</a-select-option>
          <a-select-option value="png">png</a-select-option>
          <a-select-option value="webp">webp</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="Width">
        <a-input-number
          class="item"
          v-model:value="config.width"
          :min="1"
          :max="100000"
        />
      </a-form-item>
      <a-form-item>
        <a-button type="primary" @click.prevent="onSave">保存</a-button>
        <a-button style="margin-left: 10px" @click="onReset">重置</a-button>
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

    watch(config, () => {
      console.log(config.value);
    });

    const onSave = () => {
      setConfig(config.value);
    };

    const onReset = () => {
      console.log("resetFields");
    };

    return {
      config,
      setOutDir,
      onSave,
      onReset
    };
  }
});
</script>

<style scoped>
.item {
  width: 200px;
}
</style>
