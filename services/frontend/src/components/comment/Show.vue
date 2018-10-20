<template>
<div class="container border">
  <div class="row">
    {{ comment.Body }}
  </div>
  <div class="row">
    <small>Created at: {{ comment.CreatedAt }}; Modified at:{{ comment.ModifiedAt }}</small>
  </div>
  <div class="row">
    <div class="button button-outline" @click="showCommentForm">Reply</div>
    <div class="button button-clear" style="margin-left:10px;" @click="showEditForm">Edit</div>
    <div class="button button-clear" style="margin-left:10px;" @click="deleteComment">Delete</div>
  </div> 
    <div v-show="replying">
      <submitCommentForm :postUID="comment.PostUID" :parentUID="comment.UID"></submitCommentForm>
    </div>
    <div v-show="editing">
      <editCommentForm :comment="comment"></editCommentForm>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

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
      editing: false
    }
  },
  methods: {
    deleteComment() {
      axios
        .delete('http://localhost:8081/api/posts/' + this.comment.PostUID + '/comments/' + this.comment.UID, {
          headers: {'Access-Control-Allow-Origin': '*',
          }
        })
        .catch(error => {
          console.log(error)
          this.errored = true
        })
    },
    showCommentForm() {
      this.replying = true
    },
    closeCommentForm() {
      this.replying = false
      this.$parent.fetchComments()
    },
    showEditForm() {
      this.editing = true
    },
    closeEditForm() {
      this.editing = false
      this.$parent.fetchComments()
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