module.exports = function(grunt) {
    grunt.initConfig({
        watch: {
            less: {
                files: ['style/*.less'],
                tasks: ['less']
            }
        },

        less: {
            dist: {
                options: {
                    style: 'compressed'
                },
                files: {
                    'style/style.css': 'style/importer.less'
                }
            }
        }
    });

    grunt.loadNpmTasks('grunt-contrib-watch');
    grunt.loadNpmTasks('grunt-contrib-less');
    grunt.registerTask('default', ['watch']);
};
