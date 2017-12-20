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
    formatTime: function (time) {
      var d = new Date(time * 1000);
      return d.getFullYear() + '-' + ('0'+(d.getMonth()+1)).slice(-2) + '-' + ('0'+d.getDate()).slice(-2) + ' ' +
        ('0'+d.getHours()).slice(-2) + ':' + ('0'+d.getMinutes()).slice(-2) + ':' + ('0'+d.getSeconds()).slice(-2) + '.' +
        ('00'+d.getMilliseconds()).slice(-3);
    },
    formatSize: function (size) {
      if (size > 1024) {
        return parseInt(size / 1024) + ' KB';
      } else {
        return size + ' B';
      }
    },
    latencyStr: function(latencyNs) {
      var latencyMs = latencyNs / 1000000;
      if (latencyMs > 1000) {
        return parseInt(latencyMs / 1000) + ' s';
      } else {
        return parseInt(latencyMs) + ' ms';
      }
    },
    logLevelToSymbol: function(level) {
      return level.toLowerCase().substr(0, 1);
    },
    blurFilterInput: function(e) {
        e.target.blur();
    },
  }
});

var source = new EventSource('/event/logs');
source.onmessage = function(e) {
  var log = JSON.parse(e.data);
  console.log(log);
  if (log) {
    log.expanded = false;
    app.logs.unshift(log);
  }
};
