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
      <div class="button button-outline" style="margin-left:10px;" @click="cancel">Cancel</div>
    </form>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'editPostForm',
  props: ['post'],
  data () {
    return {
      errors: [],
      title: this.post.Title,
      URL: this.post.URL
    }
  },
  methods: {
    checkForm(e) {
      this.errors = []
      if (!this.title) {
        this.errors.push("Post title can't be empty.")
      }
      if(!this.errors.length) {
        this.updatePost()
        this.post.Title = this.title
        this.post.URL = this.URL
      }
      e.preventDefault()
    },
    updatePost() {
      axios
        .patch('http://localhost:8081/api/posts/' + this.post.UID, JSON.stringify({'title': this.title, 'url': this.URL})
      )
      .then(response => {
        console.log(response)
        this.$parent.closeEditForm()
      })
      .catch(error => {
        console.log(error)
      })
    },
    cancel() {
      this.$parent.closeEditForm()
    }
  }
}
</script>
