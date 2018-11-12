import Vue from 'vue'
import Vuex from 'vuex'

import {HTTP} from '@/util/http'
import toast from '@/util/toast'

Vue.use(Vuex);

const LOGIN = "LOGIN"
const LOGIN_SUCCESS = "LOGIN_SUCCESS"
const LOGOUT = "LOGOUT"

export default new Vuex.Store({
  state: {
    isLoggedIn: localStorage.getItem("token")
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
      commit(LOGIN)
      return new Promise(resolve => {
        HTTP.post('auth', JSON.stringify({'username': payload['username'], 'password': payload['password']}))
        .then(response => {
          console.log(response.data)
          toast.success('Logged in')
          localStorage.setItem("token", response.data['Token'])
          commit(LOGIN_SUCCESS)
          resolve()
        })
        .catch(error => {
          toast.error(error.message)
        })
      });
    },
    logout({ commit }) {
      localStorage.removeItem("token");
      commit(LOGOUT);
    }
 },
 getters: {
  isLoggedIn: state => {
    return state.isLoggedIn
   }
  }
});