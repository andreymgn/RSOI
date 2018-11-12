<template>
  <div id="app">
    <div id="nav">
      <router-link :to="'/'">Home</router-link> | 
      <router-link :to="'/submit'">New post</router-link> | 
      <router-link v-if="!isLoggedIn" :to="'/login'">Log in</router-link> | 
      <a href="#" v-if="isLoggedIn" @click="logout">Logout</a>
      <router-link v-else :to="'/register'">Register</router-link>
    </div>
    <router-view/>
  </div>
</template>

<script>
import Vuex from 'vuex'
export default {
  methods: {
    ...Vuex.mapActions(["logout"])
  },
  computed: {
    ...Vuex.mapGetters(["isLoggedIn"])
  }
}
</script>


<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  background-color: #fefefe;
}
#nav {
  padding: 30px;
}

#nav a {
  font-weight: bold;
  color: #2c3e50;
}

#nav a.router-link-exact-active {
  color: #9b4dca;
}
</style>
