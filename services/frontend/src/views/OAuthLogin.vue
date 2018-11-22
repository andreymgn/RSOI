<template>
<div class="container" v-if="responseType === 'code'">
  App {{ appName }} asks for access to your account.
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
    <br>
    <input class="button-primary" type="submit" value="Log in">
  </form>
</div>
<div v-else>
  response_type != code
  </div>
</template>

<script>
import {HTTP} from '@/util/http'
import toast from '@/util/toast'

export default {
  name: 'oauthLogin',
  data () {
    return {
      errors: [],
      username: null,
      password: null,
      appName: null,
      appOwnerUsername: null,
      clientID: this.$route.query.client_id,
      redirectURI: this.$route.query.redirect_uri,
      responseType: this.$route.query.response_type
    }
  },
  created() {
    this.getAppInfo()
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
        this.getCode()
      }
      e.preventDefault()
    },
    getAppInfo() {
      HTTP.get('oauth/app/' + this.clientID)
      .then(response => {
        this.appName = response.data.Name
      })
      .catch(error => {
        toast.error(error.message)
      })
    },
    getCode() {
      HTTP.post('oauth/authorize', JSON.stringify({'AppUID': this.clientID, 'username': this.username, 'password': this.password}))
      .then(response => {
        const code = response.data.Code
        window.location = this.redirectURI + '?code=' + code
      })
      .catch(error => {
        toast.error(error.message)
      })
    }
  }
}
</script>