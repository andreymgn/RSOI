<template>
<div class="container">
  <div v-if="pending">
    Logging in
  </div>
  <form @submit="checkForm" novalidate="true">
    <div class="row error" v-if="errors.length">
      <b>Please correct the following error(s):</b>
      <ul>
        <li v-for="error in errors" :key="error">{{ error }}</li>
      </ul>
    </div>
      <label for="username">Username</label>
      <input type="text" name="username" id="username" v-model="username">
      <label for="password">Password</label>
      <input type="password" name="password" id="password" v-model="password">
      <label for="remember">Remember me</label>
      <input type="checkbox" name="remember" id="remember" value="Remember" v-model="remember">
      <br>
      <input class="button-primary" type="submit" value="Log in">
    </form>
  </div>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'loginForm',

  data () {
    return {
      errors: [],
      username: null,
      password: null,
      remember: false
    }
  },
  computed: {
    ...mapState(['pending'])
  },
  methods: {
    checkForm(e) {
      this.errors = [];
      if (!this.username) {
        this.errors.push("Username is required.")
      }
      if (!this.password) {
        this.errors.push("Password is required.")
      }
      if(!this.errors.length) {
        this.login()
      }
      e.preventDefault()
    },
    login() {
      this.$store.dispatch("login", {username: this.username, password: this.password, refresh: this.remember})
      .then(() => this.$router.push('/'))
    }
  }
}
</script>