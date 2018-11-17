import Vue from 'vue'
import Vuex from 'vuex'

import {HTTP} from '@/util/http'
import toast from '@/util/toast'

Vue.use(Vuex);

const LOGIN = 'LOGIN'
const LOGIN_SUCCESS = 'LOGIN_SUCCESS'
const LOGOUT = 'LOGOUT'

export default new Vuex.Store({
  state: {
    pending: false,
    isLoggedIn: localStorage.getItem('accessToken')
  },
  mutations: {
    [LOGIN] (state) {
      state.pending = true;
    },
    [LOGIN_SUCCESS] (state) {
      state.isLoggedIn = true;
      state.pending = false;
    },
    [LOGOUT](state) {
      state.isLoggedIn = false;
    }
  },
  actions: {
    login({ commit }, payload) {
      return new Promise(resolve => {
        commit(LOGIN)
        HTTP.post('auth/token', JSON.stringify({'username': payload['username'], 'password': payload['password'], 'refresh': payload['refresh']}))
        .then(response => {
          toast.success('Logged in')
          localStorage.setItem('accessToken', response.data['AccessToken'])
          localStorage.setItem('UID', response.data['UID'])
          if (response.data['RefreshToken']) {
            localStorage.setItem('refreshToken', response.data['RefreshToken'])
          }
          localStorage.setItem('username', payload['username'])
          resolve()
          commit(LOGIN_SUCCESS)
        })
        .catch(error => {
          toast.error(error.message)
        })
      })
    },
    logout({ commit }) {
      localStorage.removeItem('accessToken')
      localStorage.removeItem('UID')
      localStorage.removeItem('refreshToken')
      localStorage.removeItem('username')
      commit(LOGOUT);
    },
    refresh() {
      HTTP.post('auth/refresh', JSON.stringify({'token': localStorage.getItem('refreshToken')}))
      .then(response => {
        localStorage.setItem('accessToken', response.data['AccessToken'])
        if (response.data['RefreshToken']) {
          localStorage.setItem("refreshToken", response.data['RefreshToken'])
        }
      })
      .catch(error => {
        toast.error(error.message)
      })
    }
 },
 getters: {
  isLoggedIn: state => {
    return state.isLoggedIn
   }
  }
});