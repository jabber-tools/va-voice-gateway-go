// https://stackoverflow.com/questions/36290358/can-pm2-be-used-with-compiled-c-programs
module.exports = {
  apps: [
    {
      name: 'va-voice-gateway-go',
      script: '/appl/va-voice-gateway-go/va-voice-gateway /appl/va-voice-gateway-go-conf/appconf.toml',
      exec_interpreter: "none",
      exec_mode  : "fork_mode",
      kill_timeout : 40000
    }
  ]
}
