declare module "*.vue" {
  import { defineComponent } from "vue";
  const component: ReturnType<typeof defineComponent>;
  export default component;
}

interface Todos {
  SaveList(string: String): Promise<string>;
  LoadList(): Promise<string>;
  SaveAs(string: String): Promise<any>;
  LoadNewList(): Promise<string>;
}

interface FileManager {
  HandleFile(string: string): Promise<string>;
  Convert(): Promise<any>;
  OpenFile(path: string): Promise<any>;
}

interface AppConfig {
  webpOpt: { lossless: boolean; quality: number };
  pngOpt: { quality: number };
  jpegOpt: { quality: number };
  outDir: string;
  target: string;
}

interface Config {
  SetOutDir(): Promise<any>;
  GetAppConfig(): Promise<AppConfig>;
  SetConfig(config: string): Promise<any>;
  OpenOutputDir(): Promise<any>;
  Clear(): Promise<any>;
  RestoreDefaults(): Promise<any>;
}

interface Stat {
  GetStats(): Promise<any>;
}

declare interface Window {
  backend: {
    basic(): Promise<string>;
    HandleResize(): Promise<any>;
    Todos: Todos;
    FileManager: FileManager;
    Config: Config;
    Stat: Stat;
  };
}
