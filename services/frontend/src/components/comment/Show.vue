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
      editing: false
    }
  },
  methods: {
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
    },
    deleteComment() {
      var postUID = this.comment.PostUID
      var commentUID = this.comment.UID
      HTTP.delete('posts/' + postUID + '/comments/' + commentUID)
          .then(response => {
            console.log(response)
            toast.success('Comment deleted')
            this.$parent.deleteComment(postUID, commentUID)
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