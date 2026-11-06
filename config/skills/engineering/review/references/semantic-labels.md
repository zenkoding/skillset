# Semantic labels

Every finding starts with one label, so the author knows what the reviewer expects without reading tone into the text. The rest of the comment stays short and direct; the label carries the intent.

| Label | The reviewer means | The author is expected to | Blocks merge |
| --- | --- | --- | --- |
| `Crucial` | This must not merge as is: it causes harm (a security hole, data loss, a broken contract). | Change it before merge. | yes |
| `Important` | This matters and a change is expected. | Change it, or say why not and discuss. | no |
| `Suggestion` | Worth considering. | Think about it; accepting is the author's call. | no |
| `Question` | An honest question the reviewer can't answer from the code. Never rhetorical or leading: a point dressed as a question takes another label. | Answer it. | no |
| `Hint` | Useful for the future, such as a library, a pattern or a doc. | Nothing now. | no |
| `Remark` | An observation, including praise. | Nothing. | no |

How strongly a reviewer feels sets the label as much as the content does: consistent naming can be a `Suggestion` in one repo and `Important` in another. A repo's `REVIEW.md` can say which. Replies to a labelled finding don't need labels; add one only when a reply makes a new ask.

```text
Question: Where does this formula come from?
  > It's the quadratic formula for where the sweep-line parabolas meet.
Suggestion: Extract it as solveQuadratic(a, b, c), so the name says that.
```

## Source

[Semantic Code Reviews](https://www.m31coding.com/blog/semantic-reviews.html) on the M31 coding blog (2023), which builds on the comment-labelling idea in Google's engineering practices. Restated here so the skill doesn't depend on the page staying up.
