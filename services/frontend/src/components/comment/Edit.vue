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
        editComment() {
            HTTP.patch('posts/' + this.comment.PostUID + '/comments/' + this.comment.UID, JSON.stringify({'body': this.body }))
                .then(response => {
                    console.log(response)
                    toast.success('Comment changed')
                    this.$parent.closeEditForm()
                })
                .catch(error => {
                    toast.error(error.message)
                })
        },
        cancel() {
            this.$parent.closeEditForm()
        }
    }
}
</script>