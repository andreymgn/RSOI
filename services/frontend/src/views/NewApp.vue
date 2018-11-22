<template>
<div class="container">
  <form @submit="checkForm" novalidate="true">
    <div class="row error" v-if="errors.length">
      <b>Please correct the following error(s):</b>
      <ul>
        <li v-for="error in errors" :key="error">{{ error }}</li>
      </ul>
    </div>
      <label for="name">App name</label>
      <input type="text" name="name" id="name" v-model="name">
      <br>
      <input class="button-primary" type="submit" value="Create app">
    </form>
    <div v-if="success">
      App id: {{ appID }}
      App secret: {{ appSecret }}
    </div>
  </div>
</template>

<script>
import {HTTP} from '@/util/http'
import toast from '@/util/toast'

export default {
  name: 'newApp',
  data () {
    return {
      errors: [],
      name: null,
      appID: null,
      appSecret: null,
      success: false
    }
  },
  methods: {
    checkForm(e) {
      this.errors = [];
      if (!this.name) {
        this.errors.push("App name is required.")
      }
      if(!this.errors.length) {
        this.createApp()
      }
      e.preventDefault()
    },
    createApp() {
      HTTP.post('oauth/app', JSON.stringify({'name': this.name}), {headers: {'Authorization': 'Bearer ' + localStorage.getItem('accessToken')}})
      .then(response => {
        toast.success('App created')
        this.success = true
        this.appID = response.data.ID
        this.appSecret = response.data.Secret
      })
      .catch(error => {
        toast.error(error.message)
      })
    }
  }
}
</script>