package store

const (
	getMostTagChangesStats = `
SELECT count(th.user_id) as count, u.name
FROM tag_histories th
  JOIN users u
  ON th.user_id=u.id
GROUP BY th.user_id
ORDER BY count;
`

	getMostImagesUploaded = `
SELECT count(img.owner_id) as count, u.name
FROM images img
  JOIN users u
  ON img.owner_id=u.id
GROUP BY img.owner_id
ORDER BY count DESC;
`
)
