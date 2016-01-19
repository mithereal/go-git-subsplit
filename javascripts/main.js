
$(function() {
    $('#github-commits').githubInfoWidget(
        { user: 'mithereal', repo: 'go-git-subsplit', branch: 'master', last: 15, limitMessageTo: 30 });
});
