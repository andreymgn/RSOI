<template>
  <div class="container">
    <div class="row">
      <post v-if="post" :post="post"></post>
    </div>
    <div v-if="isLoggedIn" class="button float-left" @click="showCommentForm">New comment</div>
    <div class="row">
      <div v-if="editing">
        <submitCommentForm :postUID="post.UID"></submitCommentForm>
      </div>
    </div>
    <div class="column" v-if="comments && comments.length > 0">
      <comment v-for="comment in comments" :key="comment.UID" :comment="comment"></comment>
    </div>
    <div class="row" v-else>No comments yet</div>
    <button v-show="pageNumber > 0" @click="loadPrevious">&lt;</button>
    <button v-if="itemsLoaded == pageSize" @click="loadNext" style="margin-left:10px;">&gt;</button>
  </div>
</template>

<script>
import Vuex from 'vuex'

import {HTTP} from '@/util/http'
import toast from '@/util/toast'

import Post from '@/components/post/Show.vue'
import Comment from '@/components/comment/Show.vue'
import SubmitCommentForm from '@/components/comment/New.vue'

export default {
  name: 'comments',
  components: {
    Comment,
    Post,
    SubmitCommentForm
  },
  props: {
    type: String
  },
  data () {
    return {
      post: null,
      comments: null,
      editing: false,
      pageNumber: null,
      pageSize: null,
      itemsLoaded: 0
    }
  },
  computed: {
    ...Vuex.mapGetters(['isLoggedIn'])
  },
  created () {
    this.fetchData()
  },
  watch: {
    '$route': 'fetchData'
  },
  methods: {
    fetchData() {
      this.fetchPost()
      this.fetchComments(0, 10)
    },
    fetchPost() {
      HTTP.get('posts/' + this.$route.params.uid)
      .then(response => {
        this.post = response.data
      })
      .catch(error => {
        if (error.response.status == 422)
          this.$router.push('/422')
        if (error.response.status == 404)
          this.$router.push('/404')
        toast.error(error.message)
      })
    },
    fetchComments(pageNumber, pageSize) {
      HTTP.get('posts/' + this.$route.params.uid + '/comments/', {
        params: {
          size: pageSize,
          page: pageNumber
        }
      })
      .then(response => {
        this.comments = response.data.Comments
        this.itemsLoaded = this.comments.length
        this.pageNumber = response.data.PageNumber
        this.pageSize = response.data.PageSize
      })
      .catch(error => {
        toast.error(error.message)
      })
    },
    deleteComment(postUID, commentUID) {
      for (var i = 0; i < this.comments.length; i++) {
        if (this.comments[i].UID == commentUID) {
          this.$delete(this.comments, i)
        }
      }
    },
    loadPrevious() {
      this.fetchComments(this.pageNumber - 1)
    },
    loadNext() {
      this.fetchComments(this.pageNumber + 1)
    },
    showCommentForm() {
      this.editing = true
    },
    closeCommentForm() {
      this.editing = false
      this.fetchComments()
    }
  }
}
</script>

<style>
.noborder * {
  border: 0px
}
</style>
