<template>
  <div v-bind:class="{'container border': comments,  'container ': !comments}">
    <div class="row">
      <div class="column" v-if="post.NumViews != -1">
        <div class="row" style="cursor: pointer;" @click="like">&#x1F525; {{ post.NumLikes }}</div>
        <div class="row" style="cursor: pointer;" @click="dislike">&#x1F4A9; {{ post.NumDislikes }}</div>
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
        <div class="row" v-if="post.NumViews != -1">
          <i class="fa fa-eye"><small> {{ post.NumViews }}</small></i>
        </div>
      </div>
    </div>
    <div class="row">
      <small>Created {{ post.CreatedAt | timeAgo }} by {{ username }} </small>
      <small v-if="post.CreatedAt != post.ModifiedAt">; Modified: {{ post.ModifiedAt | timeAgo}}</small>
    </div>
    <div v-if="comments" class="row">
      <router-link :to="'/post/' + post.UID"><small>Read comments</small></router-link>
    </div>
    <div v-if="post.UserUID === uid" class="row">
      <div class="button" @click="showEditForm">Edit</div>
      <div class="button button-outline" style="margin-left:10px;" @click="deletePost">Delete</div>
    </div>
    <div class="row" v-if="editing">
      <editPostForm :post="post"></editPostForm>
    </div>
  </div>
</template>

<script>
import {HTTP} from '@/util/http'
import toast from '@/util/toast'

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
      uid: localStorage.getItem('UID'),
      username: ''
    }
  },
  created () {
    this.fetchUser()
  },
  methods: {
    like(retry=true) {
      HTTP.patch('posts/' + this.post.UID + '/like', '', {headers: {'Authorization': 'Bearer ' + localStorage.getItem('accessToken')}})
        .then(() => {
          this.post.NumLikes++
        })
        .catch(error => {
          if (retry && error.response.status === 403) {
            this.handle403(() => this.like(retry=false))
          }
          toast.error(error.message)
        })
    },
    dislike(retry=true) {
      HTTP.patch('posts/' + this.post.UID + '/dislike', '', {headers: {'Authorization': 'Bearer ' + localStorage.getItem('accessToken')}})
        .then(() => {
          this.post.NumDislikes++
        })
        .catch(error => {
          if (retry && error.response.status === 403) {
            this.handle403(() => this.dislike(retry=false))
          }
          toast.error(error.message)
        })
    },
    deletePost(retry=true) {
      var postUID = this.post.UID
      HTTP.delete('posts/' + this.post.UID, {headers: {'Authorization': 'Bearer ' + localStorage.getItem('accessToken')}})
        .then(() => {
          toast.success('Post deleted')
          this.$parent.deletePost(postUID)
        })
        .catch(error => {
          if (retry && error.response.status === 403) {
            this.handle403(() => this.deletePost(retry=false))
          }
          toast.error(error.message)
        })
    },
    showEditForm() {
      this.editing = true
    },
    closeEditForm() {
      this.editing = false
    },
    handle403(func) {
      if (localStorage.getItem('refreshToken')) {
        this.$store.dispatch('refresh')
        func()
      } else {
        this.$store.dispatch('logout')
        this.$router.push('/login/')
      }
    },
    fetchUser() {
      HTTP.get('user/' + this.post.UserUID)
      .then((response) => {
        this.username = response.data.Username
      })
      .catch(error => {
        toast.error(error.message)
      })
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