<template>
  <div v-bind:class="{'container border': comments,  'container ': !comments}">
    <div class="row">
      <div class="column">
        <div class="row" @click="like">+{{ post.NumLikes }}</div>
        <div class="row" @click="dislike">-{{ post.NumDislikes }}</div>
      </div>
      <div class="column column-90">
        <div class="row float-left">
          <template v-if="post.URL">
            <a :href="post.URL" target="_blank" rel="noopener">{{ post.Title }} <small><em>({{ post.URL }})</em></small></a>
          </template>
          <template v-else>
            {{ post.Title }}
          </template>
        </div>
        <div class="row">
          <small>Views: {{ post.NumViews }}</small>
        </div>
      </div>
    </div>
    <div class="row">
      <small>Created at: {{ post.CreatedAt }}; Modified at: {{ post.CreatedAt }}</small>
    </div>
    <div v-if="comments" class="row">
      <router-link :to="'/post/' + post.UID"><small>Read comments</small></router-link>
    </div>
    <div class="row">
      <div class="button" @click="showEditForm">Edit</div>
      <div class="button button-outline" style="margin-left:10px;" @click="deletePost">Delete</div>
    </div>
    <div class="row" v-if="editing">
      <editPostForm :post="post"></editPostForm>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

import EditPostForm from '@/components/post/Edit.vue'

export default {
  name: 'post',
  components: {
    EditPostForm
  },
  props: ['post', 'comments'],
  data () {
    return {
      editing: false,
    }
  },
  methods: {
    like() {
      this.post.NumLikes++
      axios
        .get('http://localhost:8081/api/posts/' + this.post.UID + '/like', {
          headers: {'Access-Control-Allow-Origin': '*',
          }
        })
        .catch(error => {
          console.log(error)
          this.errored = true
        })
    },
    dislike() {
      this.post.NumDislikes++
      axios
        .get('http://localhost:8081/api/posts/' + this.post.UID + '/dislike', {
          headers: {'Access-Control-Allow-Origin': '*',
          }
        })
        .catch(error => {
          console.log(error)
          this.errored = true
        })
    },
    deletePost() {
      axios
        .delete('http://localhost:8081/api/posts/' + this.post.UID, {
          headers: {'Access-Control-Allow-Origin': '*',
          }
        })
        .catch(error => {
          console.log(error)
          this.errored = true
        })
    },
    showEditForm() {
      this.editing = true
    },
    closeEditForm() {
      this.editing = false
    }
  }
}
</script>

<style>
  .border {
    border: 1px solid rgb(84, 34, 178);
    border-radius: 1px;
    margin-top: 2px;
    margin-bottom: 2px;
  }
</style>