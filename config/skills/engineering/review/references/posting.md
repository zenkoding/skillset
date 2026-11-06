# Posting a review

One `POST` creates the whole review, summary and inline comments, so the author gets one notification.

## Comments

- An inline comment goes only on a line inside a diff hunk on the PR head (`side: "RIGHT"`; context lines in a hunk count). Check each line against the head, not the merge ref or the diff's own numbering: `git show <head-sha>:<path> | grep -n '<code>'`.
- A range also takes `start_line` and `start_side`, with `start_line` below `line`.
- Every other finding goes in the summary body, with its `path:line`.
- One finding that covers several sites is anchored at one and names the others (`also flow.go:160`).

## Payload

Write the summary to `review-body.md` and the inline comments to `comments.json` (an array of `{path, line, side, body}`, or `[]`), both in your own scratch or worktree directory. Build the payload with `jq` so quoting in the body can't break it:

```bash
jq -n --arg sha "$HEAD_SHA" --rawfile body review-body.md --slurpfile comments comments.json \
  '{commit_id: $sha, event: "COMMENT", body: $body, comments: $comments[0]}' > review.json
gh api repos/{owner}/{repo}/pulls/<pr>/reviews --method POST --input review.json
```

The event is always `COMMENT`: a person approves and merges.

## Failures

- On HTTP 422 naming a comment location, move every inline comment into the body with its `path:line`, rebuild with `comments: []`, and retry once. Any other failure: stop and report the command and its error; never fall back to a plain PR comment.
- If the head SHA moved while you reviewed, post on the reviewed SHA and say in the body that newer commits are unreviewed.

## After posting

- Read back `gh api repos/{owner}/{repo}/pulls/<pr>/reviews` and confirm exactly one `<!-- review-pr -->` review on the head SHA.
- Remove the label: `gh api -X DELETE repos/{owner}/{repo}/issues/<pr>/labels/ready-for-agent` (a 404 means it wasn't there).
