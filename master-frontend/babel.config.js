module.exports = {
  "presets": [
    [
      "@vue/app",
      {
        "useBuiltIns": "entry"
      }
    ]
  ],
  "plugins": [
    ["@babel/proposal-decorators", { "legacy": true }],
    ["@babel/proposal-class-properties", { "loose": true }]
  ]
}
