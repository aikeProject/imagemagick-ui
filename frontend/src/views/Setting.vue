<template>
  <header class="p-3 flex justify-end items-center">
    <router-link to="/" class="flex item-center text-blue-500 text-xl">
      <RollbackOutlined />
    </router-link>
  </header>
  <main class="my-2 px-4">
    <a-form
      class="space-y-3"
      :model="config"
      :label-col="labelCol"
      :wrapper-col="wrapperCol"
    >
      <a-form-item label="文件目录" class="mb-0">
        <div
          @click="setOutDir"
          style="min-width: 200px;height: 32px;font-size: 15px;"
          class="flex items-center flex-auto bg-gray-100 rounded py-1 px-3 text-gray-400 tracking-wider font-medium cursor-pointer hover:bg-gray-200 hover:text-gray-500"
        >
          {{ config.outDir }}
        </div>
      </a-form-item>
      <a-form-item label="目标文件类型">
        <a-select v-model:value="config.target" placeholder="选择文件类型">
          <a-select-option
            :key="item"
            v-for="item in covertType"
            :value="'.' + item"
          >
            {{ item }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="width(宽) x height(高)">
        <div class="space-x-5">
          <a-input-number
            v-model:value="config.width"
            :min="0"
            :max="100000"
            :step="10"
          />
          <a-input-number
            v-model:value="config.height"
            :min="0"
            :max="100000"
            :step="10"
          />
        </div>
      </a-form-item>
      <a-form-item label="gif:delay(延迟)">
        <a-input-number
          v-model:value="config.delay"
          :min="0"
          :max="100000"
          :step="1"
        />
      </a-form-item>
      <a-form-item label="resolution(分辨率)">
        <a-input-number
          v-model:value="config.resolution"
          :min="0"
          :max="100000"
          :step="1"
        />
      </a-form-item>
      <a-form-item label="sharpen(锐化)">
        <a-input-number
          v-model:value="config.sharpen"
          :min="0"
          :max="100000"
          :step="1"
        />
      </a-form-item>
      <a-form-item :wrapper-col="{ span: 19, offset: 5 }">
        <a-button type="primary" @click.prevent="onSave">保存</a-button>
        <a-button style="margin-left: 10px" @click="onReset">重置</a-button>
      </a-form-item>
    </a-form>
  </main>
</template>

<script lang="ts">
import { computed, defineComponent, watch } from "vue";
import { message } from "ant-design-vue";
import { RollbackOutlined } from "@ant-design/icons-vue";
import { useStore } from "vuex";
import { useDebounceFn } from "composables/useDebounceFn";

export default defineComponent({
  name: "Setting",
  components: { RollbackOutlined },
  setup() {
    const covertType = ["jpg", "png", "webp", "gif", "pdf"];
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
      $store
        .dispatch("setConfig", v)
        .then(() => {
          message.success("保存成功");
        })
        .catch(() => {
          message.error("保存失败");
        });
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
      covertType,
      config,
      labelCol: { span: 5 },
      wrapperCol: { span: 19 },
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
