// https://stackoverflow.com/questions/36290358/can-pm2-be-used-with-compiled-c-programs
module.exports = {
  apps: [
    {
      name: 'va-voice-gateway-go',
      script: '/appl/va-voice-gateway-go/va-voice-gateway /appl/va-voice-gateway-go-conf/appconf.toml',
      exec_interpreter: "none",
      exec_mode  : "fork_mode",
      kill_timeout : 40000,
      env: {
        loglevel_asterisk: "debug",
        loglevel_asteriskclient: "debug",
        loglevel_config: "debug",
        loglevel_gateway: "debug",
        loglevel_nlp: "debug",
        loglevel_nlpactor: "debug",
        loglevel_stt: "debug",
        loglevel_sttactor: "debug",
        loglevel_tts: "debug",
        loglevel_utils: "debug",
        loglevel_main: "debug"
      }
    }
  ]
}
