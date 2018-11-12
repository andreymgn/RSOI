import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store';

import TimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
TimeAgo.locale(en)

const timeAgo = new TimeAgo('en-US')

Vue.filter('timeAgo', function (time) {
  return timeAgo.format(Date.parse(time))
})

Vue.config.productionTip = false

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
