<template>
<div class="container border">
  <div class="row">
    By {{ username }}:
  </div>
  <div class="row">
    {{ comment.Body }}
  </div>
  <div class="row">
    <small>Created {{ comment.CreatedAt | timeAgo }}</small>
    <small v-if="comment.CreatedAt != comment.ModifiedAt">; Modified: {{ comment.ModifiedAt | timeAgo}}</small>
  </div>
  <div class="row">
    <div class="button button-outline" @click="loadReplies">Load replies</div>
    <div v-if="comment.UserUID === uid" class="button button-outline" style="margin-left:10px;" @click="showCommentForm">Reply</div>
    <div v-if="comment.UserUID === uid" class="button button-clear" style="margin-left:10px;" @click="showEditForm">Edit</div>
    <div v-if="comment.UserUID === uid" class="button button-clear" style="margin-left:10px;" @click="deleteComment">Delete</div>
  </div> 
  <div v-show="replying">
    <submitCommentForm :postUID="comment.PostUID" :parentUID="comment.UID"></submitCommentForm>
  </div>
  <div v-show="editing">
    <editCommentForm :comment="comment"></editCommentForm>
  </div>
  <comment v-if="children" v-for="child in children" :key="child.UID" :comment="child"></comment>
  </div>
</template>

<script>
import {HTTP} from '@/util/http'
import toast from '@/util/toast'

import SubmitCommentForm from '@/components/comment/New.vue'
import EditCommentForm from '@/components/comment/Edit.vue'

export default {
  name: 'comment',
  components: {
    SubmitCommentForm,
    EditCommentForm
  },
  props: ['comment'],
  data() {
    return {
      replying: false,
      editing: false,
      children: null,
      username: '',
      uid: localStorage.getItem('UID')
    }
  },
  created () {
    this.fetchUser()
  },
  methods: {
    showCommentForm() {
      this.replying = true
    },
    closeCommentForm(cancelled) {
      this.replying = false
      if (!cancelled) {
        if (this.$parent.name == this.name) {
          this.loadReplies()
        } else {
          this.$parent.fetchComments()
        }
      }
    },
    showEditForm() {
      this.editing = true
    },
    closeEditForm() {
      this.editing = false
      this.$parent.fetchComments()
    },
    deleteComment(retry=true) {
      var postUID = this.comment.PostUID
      var commentUID = this.comment.UID
      HTTP.delete('posts/' + postUID + '/comments/' + commentUID, {headers: {'Authorization': 'Bearer ' + localStorage.getItem('accessToken')}})
        .then(() => {
          toast.success('Comment deleted')
          this.comment.Body = '[deleted]'
          this.username = '[deleted]'
        })
        .catch(error => {
          if (retry && error.response.status === 403) {
            if (localStorage.getItem('refreshToken')) {
              this.$store.dispatch('refresh')
              this.deleteComment(retry=false)
            } else {
              this.$store.dispatch('logout')
              this.$router.push('/login/')
            }
          }
          toast.error(error.message)
        })
    },
    loadReplies() {
      HTTP.get('posts/' + this.comment.PostUID + '/comments/' + this.comment.UID)
        .then(response => {
          this.children = response.data.Comments
        })
        .catch(error => {
          toast.error(error.message)
        })
    },
    fetchUser() {
      if (this.comment.UserUID === '') {
        this.username = '[deleted]'
      } else {
        HTTP.get('user/' + this.comment.UserUID)
        .then((response) => {
          this.username = response.data.Username
        })
        .catch(error => {
          toast.error(error.message)
        })
      }
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