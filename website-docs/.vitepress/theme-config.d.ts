import 'vitepress'

declare module 'vitepress/dist/client/theme-default/config' {
  export interface ThemeConfig {
    /** 文档站点展示的 SemiClaw 发布版本（来自仓库根 VERSION） */
    semiclawVersion?: string
  }
}
