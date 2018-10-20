<template>
<div class="container">
  <form @submit="checkForm" novalidate="true">
    <div class="row error" v-if="errors.length">
      <b>Please correct the following error(s):</b>
      <ul>
        <li v-for="error in errors" :key="error">{{ error }}</li>
      </ul>
    </div>
      <label for="title">Post title</label>
      <input type="text" name="title" id="title" v-model="title">
      <label for="URL">URL</label>
      <input type="URL" name="URL" id="url" v-model="URL">
      <br>
      <input class="button-primary" type="submit" value="Submit">
    </form>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'newPostForm',

  data () {
    return {
      errors: [],
      title: null,
      URL: null
    }
  },
  methods: {
    checkForm(e) {
      this.errors = [];
      if (!this.title) {
        this.errors.push("Post title is required.")
      }
      if(!this.errors.length) {
        this.submitPost()
      }
      e.preventDefault()
    },
    submitPost() {
      axios
        .post('http://localhost:8081/api/posts/',
          JSON.stringify({'title': this.title, 'url': this.URL})
      )
      .then(response => {
        console.log(response)
        this.$router.push('/post/' + response.data.UID)
      })
      .catch(error => {
        console.log(error)
        this.errored = true
      })
    }
  }
}
</script>