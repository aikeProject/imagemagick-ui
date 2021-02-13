import { createStore } from "vuex";

type ConfigKey =
  | "outDir"
  | "target"
  | "prefix"
  | "suffix"
  | "jpegOpt"
  | "pngOpt"
  | "webpOpt";

interface ConfigProp {
  key: ConfigKey;
  value: any;
}

export default createStore({
  state: {
    config: {
      outDir: "",
      target: "",
      prefix: "",
      suffix: "",
      jpegOpt: { quality: 0 },
      pngOpt: { quality: 0 },
      webpOpt: { lossless: false, quality: 0 }
    },
    stats: {
      byteCount: 0,
      imageCount: 0
    }
  },
  getters: {
    config(state) {
      return state.config;
    },
    stats(state) {
      return state.stats;
    }
  },
  mutations: {
    setConfig(state, config) {
      state.config = config;
    },
    setStats(state, s) {
      state.stats = s;
    },
    setConfigProp(state, payload: ConfigProp) {
      state.config[payload.key] = payload.value;
    },
    toggleWebpLossless(state) {
      state.config.webpOpt.lossless = !state.config.webpOpt.lossless;
    }
  },
  actions: {
    getConfig(context) {
      window.backend.Config.GetAppConfig()
        .then(config => {
          console.log("config", config);
          context.commit("setConfig", {
            ...config,
            jpegOpt: config.jpegOpt || { quality: 0 },
            pngOpt: config.pngOpt || { quality: 0 },
            webpOpt: config.webpOpt || { lossless: false, quality: 0 }
          });
        })
        .catch(err => console.error(err));
    },
    setConfig(context, config) {
      window.backend.Config.SetConfig(JSON.stringify(config))
        .then(() => {
          context.dispatch("getConfig").then();
        })
        .catch(console.error);
    },
    getStats(context) {
      window.backend.Stat.GetStats()
        .then(s => {
          context.commit("setStats", s);
        })
        .catch(err => {
          console.error(err);
        });
    },
    setStats(context, s) {
      context.commit("setStats", s);
    },
    setConfigProp(context, payload: ConfigProp) {
      context.commit("setConfigProp", payload);
    },
    toggleWebpLossless(context) {
      context.commit("toggleWebpLossless");
    }
  },
  modules: {}
});
