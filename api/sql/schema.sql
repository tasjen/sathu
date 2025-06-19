CREATE table temples (
  id UUID PRIMARY KEY,
  name_th TEXT NOT NULL,
  name_en TEXT NOT NULL,
  location POINT NOT NULL,
  address_th TEXT,
  address_en TEXT,
  contact_phone TEXT,
  founded_on DATE
);

CREATE table temple_images (
  id UUID PRIMARY KEY,
  temple_id UUID NOT NULL REFERENCES temples(id) ON DELETE CASCADE,
  image_url TEXT NOT NULL,
  caption TEXT
);

CREATE table temple_events (
  id UUID PRIMARY KEY,
  temple_id UUID NOT NULL REFERENCES temples(id) ON DELETE CASCADE,
  event_name TEXT NOT NULL,
  event_date DATE NOT NULL,
  description TEXT
);

CREATE table temple_staff (
  id UUID PRIMARY KEY,
  temple_id UUID NOT NULL REFERENCES temples(id) ON DELETE CASCADE,
  staff_name TEXT NOT NULL,
  position TEXT NOT NULL,
  contact_info TEXT
);

CREATE table temple_donations (
  id UUID PRIMARY KEY,
  temple_id UUID NOT NULL REFERENCES temples(id) ON DELETE CASCADE,
  donor_name TEXT NOT NULL,
  donation_amount DECIMAL(10, 2) NOT NULL,
  donation_date DATE NOT NULL,
  purpose TEXT
);

CREATE table temple_news (
  id UUID PRIMARY KEY,
  temple_id UUID NOT NULL REFERENCES temples(id) ON DELETE CASCADE,
  news_title TEXT NOT NULL,
  news_content TEXT NOT NULL,
  news_date DATE NOT NULL
);

CREATE table temple_reviews (
  id UUID PRIMARY KEY,
  temple_id UUID NOT NULL REFERENCES temples(id) ON DELETE CASCADE,
  review_text TEXT NOT NULL,
  review_rating INTEGER CHECK (review_rating >= 1 AND review_rating <= 5),
  review_date DATE NOT NULL
);

CREATE table temple_schedules (
  id UUID PRIMARY KEY,
  temple_id UUID NOT NULL REFERENCES temples(id) ON DELETE CASCADE,
  schedule_date DATE NOT NULL,
  schedule_time TIME NOT NULL,
  activity TEXT NOT NULL
);

CREATE table temple_sponsorships (
  id UUID PRIMARY KEY,
  temple_id UUID NOT NULL REFERENCES temples(id) ON DELETE CASCADE,
  sponsor_name TEXT NOT NULL,
  sponsorship_amount DECIMAL(10, 2) NOT NULL,
  sponsorship_date DATE NOT NULL,
  purpose TEXT
);

CREATE table temple_memberships (
  id UUID PRIMARY KEY,
  temple_id UUID NOT NULL REFERENCES temples(id) ON DELETE CASCADE,
  member_name TEXT NOT NULL,
  membership_date DATE NOT NULL,
  membership_type TEXT NOT NULL
);