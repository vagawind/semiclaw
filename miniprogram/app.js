App({
  onLaunch() {
    const settings = wx.getStorageSync("semiclaw_settings");
    if (!settings) {
      wx.setStorageSync("semiclaw_settings", {
        baseUrl: "http://localhost:8080",
        apiKey: "",
        selectedKnowledgeBaseId: ""
      });
    }
  }
});
