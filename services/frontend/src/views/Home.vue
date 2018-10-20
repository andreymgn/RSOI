<template>
  <div class="home">
    <div v-for="post in posts" :key="post.UID">
      <div v-if="post">
        <post :post="post" :comments="true"></post>
      </div>
    </div>
    <button v-show="pageNumber > 0" @click="loadPrevious">&lt;</button>
    <button v-if="posts && posts.length == pageSize" @click="loadNext" style="margin-left:10px;">&gt;</button>
  </div> 
</template>

<script>
import {HTTP} from '@/util/http'
import toast from '@/util/toast'

import Post from '@/components/post/Show.vue'

export default {
  name: 'home',
  components: {
    Post
  },

  props: {
    type: String
  },

  data () {
    return {
      posts: null,
      pageNumber: null,
      pageSize: null
    }
  },

  mounted () {
    this.getPage(0, 10)
  },

  methods: {
    getPage(pageNumber, pageSize) {
      HTTP.get('posts/', {
          params: {
            size: pageSize,
            page: pageNumber
          }
        })
        .then(response => {
          this.posts = response.data.Posts
          this.pageNumber = response.data.PageNumber
          this.pageSize = response.data.PageSize
        })
        .catch(error => {
          toast.error(error.message)
        })
    },
    loadPrevious() {
     this.getPage(this.pageNumber - 1)
    },
    loadNext() {
      this.getPage(this.pageNumber + 1)
    }
  }
}
</script>
