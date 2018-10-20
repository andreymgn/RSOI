<template>
<div class="container border">
  <div class="row">
    {{ comment.Body }}
  </div>
  <div class="row">
    <small>Created {{ comment.CreatedAt | timeAgo }}</small>
    <small v-if="comment.CreatedAt != comment.ModifiedAt">; Modified: {{ comment.ModifiedAt | timeAgo}}</small>
  </div>
  <div class="row">
    <div class="button button-outline" @click="loadReplies">Load replies</div>
    <div class="button button-outline" style="margin-left:10px;" @click="showCommentForm">Reply</div>
    <div class="button button-clear" style="margin-left:10px;" @click="showEditForm">Edit</div>
    <div class="button button-clear" style="margin-left:10px;" @click="deleteComment">Delete</div>
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
      children: null
    }
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
    deleteComment() {
      var postUID = this.comment.PostUID
      var commentUID = this.comment.UID
      HTTP.delete('posts/' + postUID + '/comments/' + commentUID)
          .then(response => {
            console.log(response)
            toast.success('Comment deleted')
            // do not delete parent comment if child gets deleted
            if (this.$parent.name != this.name) {
              this.$parent.deleteComment(postUID, commentUID)
            } else {
              for (var i = 0; i < this.$parent.children.length; i++) {
                if (this.$parent.children[i].UID == commentUID) {
                  this.$parent.$delete(this.$parent.children, i)
                }
              }
            }
          })
          .catch(error => {
            toast.error(error.message)
          })
    },
    loadReplies() {
      HTTP.get('posts/' + this.comment.PostUID + '/comments/' + this.comment.UID)
        .then(response => {
          console.log(response)
          this.children = response.data.Comments
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