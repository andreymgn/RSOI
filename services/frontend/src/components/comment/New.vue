<template>
    <div class="container">
        <form @submit="checkForm" novalidate="true">
            <div class="row error" v-if="errors.length">
                <b>Please correct the following error(s):</b>
                <ul>
                    <li v-for="error in errors" :key="error">{{ error }}</li>
                </ul>
            </div>
            <label for="body">Comment body</label>
            <textarea type="body" name="body" id="body" v-model="body"></textarea>
            <br>
            <input class="button-primary" type="submit" value="Submit">
            <div class="button button-outline" style="margin-left:10px;" @click="cancel">Cancel</div>
        </form>
    </div>
</template>

<script>
import {HTTP} from '@/util/http'
import toast from '@/util/toast'

export default {
    name: 'newCommentForm',
    props: ['postUID', 'parentUID'],
    data () {
        return {
            errors: [],
            body: null
        }
    },
    methods: {
        checkForm(e) {
            this.errors = [];
            if (!this.body) {
                this.errors.push("Comment body can't be empty.")
            }
            if (!this.errors.length) {
                this.submitComment()
            }
            e.preventDefault()
        },
        submitComment(retry=true) {
            HTTP.post('posts/' + this.postUID + '/comments/', JSON.stringify({'body': this.body, 'parent_uid': this.parentUID }), {headers: {'Authorization': 'Bearer ' + localStorage.getItem('accessToken')}})
            .then(() => {
                toast.success('Comment created')
                this.$parent.closeCommentForm(false)
            })
            .catch(error => {
                if (retry && error.response.status === 403) {
                    if (localStorage.getItem('refreshToken')) {
                        this.$store.dispatch('refresh')
                        this.submitComment(retry=false)
                    } else {
                        this.$store.dispatch('logout')
                        this.$router.push('/login/')
                    }
                }
                toast.error(error.message)
            })
        },
        cancel() {
            this.$parent.closeCommentForm(true)
        }
    }
}
</script>