-- unused, kept for reference, refer topicrepository for correct one
-- name: GetTopicByIDWithVideos :one
SELECT 
  t.*,
  COALESCE(
    JSON_AGG(
      JSON_BUILD_OBJECT(
        'id', yv.id,
        'video_id', yv.video_id,
        'title', yv.title,
        'video_date', yv.video_date,
        'channel_title', yv.channel_title,
        'thumbnail_url', yv.thumbnail_url,
        'expiry_at', yv.expiry_at,
        'updated_at', yv.updated_at,
        'created_at', yv.created_at,
        'updated_by', yv.updated_by
      )
    ) FILTER (WHERE yv.id IS NOT NULL),
    '[]'::json
  ) AS videos
FROM 
  topic t
LEFT JOIN 
  youtube_video yv ON t.id = yv.topic_id
WHERE 
  t.id = $1
GROUP BY 
  t.id;