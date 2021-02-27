import { createStore } from "vuex";

type ConfigKey = "outDir" | "target";

interface ConfigProp {
  key: ConfigKey;
  value: any;
}

export default createStore({
  state: {
    config: {
      outDir: "",
      target: "",
      width: 0,
      height: 0,
      delay: 0,
      resolution: 0,
      sharpen: 0
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
    }
  },
  actions: {
    getConfig(context) {
      window.backend.Config.GetAppConfig()
        .then(config => {
          context.commit("setConfig", config);
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
