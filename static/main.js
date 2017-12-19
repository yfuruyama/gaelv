var app = new Vue({
  el: '#app',
  delimiters: ['${', '}'],
  data: {
    filterText: "",
    logs: [],
  },
  watch: {
  },
  computed: {
    filteredLogs: function() {
      var re = new RegExp(this.filterText);
      return this.logs.filter(function(log) {
        return log.resource.match(re);
      });
    },
  },
  methods: {
    formattedTime: function (time) {
      var d = new Date(time * 1000);
      // TODO: zero padding
      return `${d.getFullYear()}-${(d.getMonth() + 1)}-${d.getDate()} ${d.getHours()}:${d.getMinutes()}:${d.getSeconds()}.${d.getMilliseconds()}`
           ;
    },
    latencyStr: function(latencyNs) {
      var latencyMs = latencyNs / 1000000;
      if (latencyMs > 1000) {
        return `${latencyMs / 1000} s`;
      } else {
        return `${latencyMs} ms`;
      }
    },
    toggleExpansion: function(log) {
      console.log("before: " + log.expanded);
      log.expanded = !log.expanded;
      console.log("after: " + log.expanded);
    },
    logLevelToSymbol: function(level) {
      return level.toLowerCase().substr(0, 1);
    },
  }
});

var source = new EventSource('/event/logs');
source.onmessage = function(e) {
  var log = JSON.parse(e.data);
  console.log(log);
  app.logs.push(log);
};
