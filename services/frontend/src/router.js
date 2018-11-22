import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'
import Comments from './views/Comments.vue'
import NewPost from './views/NewPost.vue'
import Register from './views/Register.vue'
import Login from './views/Login.vue'
import NewApp from './views/NewApp.vue'
import OAuthLogin from './views/OAuthLogin.vue'

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
    },
    {
      path: '/register',
      name: 'register',
      component: Register
    },
    {
      path: '/login',
      name: 'login',
      component: Login
    },
    {
      path: '/oauth/register',
      name: 'newApp',
      component: NewApp
    },
    {
      path: '/oauth/authorize',
      name: 'oauthLogin',
      component: OAuthLogin
    },
    {
      path: '/422',
      name: 'UnprocessableEntity',
      component: () => import(/* webpackChunkName: "UnprocessableEntity" */ './views/UnprocessableEntity.vue'),
    },
    {
      path: '/404',
      name: 'NotFound',
      component: () => import(/* webpackChunkName: "NotFound" */ './views/NotFound.vue'),
    }
  ]
})
