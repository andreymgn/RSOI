<template>
    <div class="container">
        <form @submit="checkForm" novalidate="true">
            <div class="row error" v-if="errors.length">
                <b>Please correct the following error(s):</b>
                <ul>
                    <li v-for="error in errors" :key="error">{{ error }}</li>
                </ul>
            </div>
            <label for="body">Edit comment</label>
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
    name: 'editCommentForm',
    props: ['comment'],
    data () {
        return {
            errors: [],
            body: this.comment.Body
        }
    },
    methods: {
        checkForm(e) {
            this.errors = [];
            if (!this.body) {
                this.errors.push("Comment body can't be empty.")
            }
            if (!this.errors.length) {
                this.editComment()
            }
            e.preventDefault()
        },
        editComment(retry=true) {
            HTTP.patch('posts/' + this.comment.PostUID + '/comments/' + this.comment.UID, JSON.stringify({'body': this.body }), {headers: {'Authorization': 'Bearer ' + localStorage.getItem('accessToken')}})
            .then(() => {
                toast.success('Comment changed')
                this.$parent.closeEditForm()
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
            this.$parent.closeEditForm()
        }
    }
}
</script>