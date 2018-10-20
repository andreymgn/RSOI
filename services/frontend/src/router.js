import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'
import Comments from './views/Comments.vue'
import NewPost from './views/NewPost.vue'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/post/:uid',
      name: 'post',
      component: Comments
    },
    {
      path: '/submit',
      name: 'submit',
      component: NewPost
    }
  ]
})
