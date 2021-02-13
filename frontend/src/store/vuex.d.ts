/**
 * @author 成雨
 * @date 2021/1/30
 * @Description:
 */

import { ComponentCustomProperties } from "vue";
import { Store } from "vuex";

declare module "@vue/runtime-core" {
  // declare your own store states
  interface State {
    config: AppConfig;
  }

  // provide typings for `this.$store`
  interface ComponentCustomProperties {
    $store: Store<State>;
  }
}
