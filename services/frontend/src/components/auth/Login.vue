<template>
<div class="container">
  <form @submit="checkForm" novalidate="true">
    <div class="row error" v-if="errors.length">
      <b>Please correct the following error(s):</b>
      <ul>
        <li v-for="error in errors" :key="error">{{ error }}</li>
      </ul>
    </div>
      <label for="username">Username</label>
      <input type="text" name="username" id="username" v-model="username">
      <label for="Password">Password</label>
      <input type="password" name="password" id="password" v-model="password">
      <br>
      <input class="button-primary" type="submit" value="Log in">
    </form>
  </div>
</template>

<script>
export default {
  name: 'loginForm',

  data () {
    return {
      errors: [],
      username: null,
      password: null
    }
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
      this.$store.dispatch("login", {username: this.username, password: this.password})
      .then(this.$router.push('/'))
    }
  }
}
</script>