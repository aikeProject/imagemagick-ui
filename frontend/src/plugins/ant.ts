/**
 * @author 成雨
 * @date 2021/2/21
 * @Description:
 */
import { App } from "vue";
import Antd from "ant-design-vue";
import "ant-design-vue/dist/antd.css";

export default (app: App) => {
  app.use(Antd);
};
