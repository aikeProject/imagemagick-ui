<template>
  <div class="flex h-screen w-screen">
    <div class="p-3 space-y-2">
      <div>
        <a-button @click="handleConvert" shape="circle">
          <ThunderboltOutlined style="vertical-align: 0;display: block;" />
        </a-button>
      </div>
      <div>
        <a-button @click="handleClear" shape="circle">
          <ClearOutlined style="vertical-align: 0;display: block;" />
        </a-button>
      </div>
      <div>
        <router-link to="setting">
          <a-button shape="circle">
            <SettingOutlined style="vertical-align: 0;display: block;" />
          </a-button>
        </router-link>
      </div>
    </div>
    <main class="flex-grow h-full flex overflow-y-auto">
      <drag-file
        v-model:show="dragShow"
        class="flex-grow my-3 mr-3"
        @change="dragChange"
      ></drag-file>
      <figure
        v-show="!!filesView.length"
        class="flex flex-wrap flex-col flex-grow self-start py-3 pr-3"
      >
        <div
          class="relative p-2 h-28 bg-white rounded-md shadow-sm flex flex-grow cursor-pointer mb-3 border border-gray-200"
          v-for="item in filesView"
          :key="item.name"
        >
          <img class="block object-cover rounded w-24" :src="item.src" alt="" />
          <div class="relative space-y-1 pl-3 flex-grow text-gray-500">
            <div class="w-80 text-sm truncate">
              文件名: {{ fName(item.name) }}
            </div>
            <div class="text-sm">大小: {{ fSize(item.size) }}</div>
            <div class="text-sm">类型: {{ item.ext }}</div>
            <div
              v-show="item.show"
              class="absolute inset-x-0 bottom-0 pl-3 pr-2"
            >
              <a-progress
                :percent="item.progress"
                :status="item.statusStr"
                size="small"
              ></a-progress>
            </div>
          </div>
          <div class="absolute top-0 right-0 p-2">
            <a-button @click="handleConvertItem(item.id)" shape="circle">
              <ThunderboltOutlined style="vertical-align: 0;display: block;" />
            </a-button>
          </div>
        </div>
      </figure>
    </main>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, watch, computed, onMounted } from "vue";
import { message, notification } from "ant-design-vue";
import {
  ThunderboltOutlined,
  ClearOutlined,
  SettingOutlined
} from "@ant-design/icons-vue";
import Wails from "@wailsapp/runtime";
import DragFile from "components/DragFile.vue";
import { fExt, readAsDataURL, fName, fSize } from "lib/file";
import { Complete, FileData } from "views/Home";
import { FileStatus } from "common/enum";
import { EventFps } from "lib/event-fps";
import useDrag from "composables/useDrag";
export default defineComponent({
  name: "Home",
  components: {
    DragFile,
    ClearOutlined,
    SettingOutlined,
    ThunderboltOutlined
  },
  setup() {
    const dragRef = ref(document.body);
    const dragShow = ref(true);
    const filesData = ref<FileData[]>([]);
    const fileSpeed = ref(1000 * 1000);
    const fileTimeMap = ref<{ [key: string]: number }>({});
    const { files } = useDrag(dragRef);

    // 根据id查询文件实例
    const getFileById = (id: string) => {
      if (!filesData.value.length) return;
      return filesData.value.find(v => v.id === id);
    };

    /**
     * 根据文件名和大小创建文件id
     * @param name
     * @param size
     */
    const createFileId = (name: string, size: number) => {
      return `${name}-${size.toString()}`;
    };

    /**
     * 初始化FileData数据
     * @param v {File}
     */
    const newFileData = (v: File): FileData => ({
      id: createFileId(v.name, v.size),
      ext: fExt(v.name),
      status: FileStatus.NotStarted,
      name: v.name,
      size: v.size,
      progress: 0,
      src: ""
    });

    // 向golang程序发送文件数据
    const sendFile = async (files: FileData[]) => {
      files = files.filter(v => v.status === FileStatus.NotStarted);
      for (const v of files) {
        // 数据已发送至golang处理程序 无需重复发送
        if (v.status === FileStatus.SendSuccess) continue;
        const timeStart = new Date().getTime();
        // 开始发送
        v.status = FileStatus.Start;
        try {
          await window.backend.Manager.HandleFile(
            JSON.stringify({
              id: v.id,
              name: v.name,
              size: v.size,
              data: v.src.split(",")[1],
              status: FileStatus.Start
            })
          );
        } catch (err) {
          notification.error({
            message: "convert",
            description: err
          });
        }
        // 发送完成
        v.status = FileStatus.SendSuccess;
        v.progress = 100;
        const timeEnd = new Date().getTime();
        fileSpeed.value = (v.size / (timeEnd - timeStart)) * 1000;
        console.log(
          "file %s \ntime %d ms \nsize => %d byte \nspeed => %fkb/s",
          v.name,
          timeEnd - timeStart,
          v.size,
          Number((fileSpeed.value / 1024).toFixed(2))
        );
      }
    };

    // 将文件转换为base64字符串
    const covertFile = async (v: File) => {
      const file = newFileData(v);
      try {
        file.src = await readAsDataURL(v);
        return file;
      } catch (e) {
        notification.error({
          message: v.name,
          description: "转换base64失败"
        });
      }
      return file;
    };

    // 构造FileData数据
    const covertFileData = async (fs: FileList): Promise<FileData[]> => {
      const files: File[] = [].slice.apply(fs);
      if (!files.length) return [];
      return (await Promise.all(files.map(v => covertFile(v)))).filter(
        v => v.src
      );
    };

    // 检查文件状态
    // 是否有文件 正在发送中或处理中
    const checkSend = () => {
      if (!filesData.value.length) return false;
      return filesData.value.some(
        v => v.status === FileStatus.Start || v.status === FileStatus.Running
      );
    };

    // 拖拽选择文件，并将文件发送至golang程序
    const dragChange = async (fs: FileList) => {
      if (!fs) return;
      const files = await covertFileData(fs);
      filesData.value = [...filesData.value, ...files];
      await sendFile(filesData.value);
    };

    // 调用golang程序处理文件
    const handleConvert = async () => {
      fileSpeed.value = 1000 * 1000;
      if (checkSend()) {
        message.warning("等待...");
        return;
      }
      // 改变文件状态
      filesData.value.forEach(v => {
        v.status = FileStatus.Running;
        v.progress = 0;
        v.statusStr = "active";
      });
      fileTimeMap.value = {};
      const { Convert } = window.backend.Manager;
      try {
        await Convert(JSON.stringify([]));
      } catch (err) {
        notification.error({
          message: "convert",
          description: err
        });
      }
    };

    // 对单个文件进行处理
    const handleConvertItem = async (id: string) => {
      const item = getFileById(id);
      if (item) {
        if (
          item.status === FileStatus.Start ||
          item.status === FileStatus.Running
        ) {
          message.warning("等待...");
          return;
        }
        fileSpeed.value = 1000 * 1000;
        item.status = FileStatus.Running;
        item.progress = 0;
        item.statusStr = "active";
        fileTimeMap.value[item.id] = 0;
        const { Convert } = window.backend.Manager;
        try {
          await Convert(JSON.stringify([item.id]));
        } catch (err) {
          notification.error({
            message: "convert",
            description: err
          });
        }
      }
    };

    // 清空
    const handleClear = async () => {
      if (checkSend()) {
        message.warning("等待...");
        return;
      }
      try {
        const { Clear } = window.backend.Manager;
        filesData.value = [];
        fileTimeMap.value = {};
        await Clear();
      } catch (err) {
        notification.error({
          message: "convert",
          description: err
        });
      }
    };

    watch(filesData, function(v) {
      // 清空操作之后，显示拖拽区域
      if (!v.length) dragShow.value = true;
    });

    watch(files, function() {
      if (checkSend()) {
        message.warning("等待...");
        return;
      }
      // 继续添加文件
      files.value && dragChange(files.value);
    });

    const filesView = computed(() => {
      return filesData.value.map(v => ({
        ...v,
        // show 进度条是否显示
        show:
          v.status === FileStatus.Done ||
          v.status === FileStatus.Error ||
          (v.progress < 100 && v.progress > 0)
      }));
    });

    onMounted(() => {
      let time = 0;
      let last = 0;
      // fps
      EventFps.on<number>("update", function(f) {
        if (filesData.value.every(v => v.progress >= 100)) return;
        // 所有文件都已经传输至golang
        if (filesData.value.every(v => v.status === FileStatus.SendSuccess)) {
          return;
        }
        // 所有文件都经由golang处理完毕
        if (filesData.value.every(v => v.status === FileStatus.Done)) {
          return;
        }
        f = f || 0;
        time += f;
        if ((time - last) * 1000 > 100) {
          for (const v of filesData.value) {
            if (
              v.status === FileStatus.Start ||
              v.status === FileStatus.Running
            ) {
              fileTimeMap.value[v.id] = (fileTimeMap.value[v.id] || 0) + f;
              v.progress = parseFloat(
                Math.min(
                  ((fileSpeed.value * fileTimeMap.value[v.id]) / v.size) * 100,
                  99
                ).toFixed(1)
              );
            }
          }
          last = time;
        }
      });

      // file:complete events 文件处理完成后收到的数据
      Wails.Events.On("file:complete", (data: Complete) => {
        console.log("file:complete");
        console.log(data);
        if (!data) return;
        const file = getFileById(data.id);
        if (file) {
          // 更新文件状态
          file.status = data.status;
          switch (data.status) {
            case FileStatus.Error:
              file.statusStr = "exception";
              break;
            case FileStatus.Done:
              file.statusStr = "success";
              break;
          }
          file.progress = 100;
        }
      });
    });

    return {
      filesView,
      dragShow,
      dragChange,
      handleConvert,
      handleConvertItem,
      handleClear,
      fName,
      fSize
    };
  }
});
</script>

<style scoped lang="stylus"></style>
