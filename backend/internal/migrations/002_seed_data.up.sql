-- Seed data: generate 1000 subscription records
-- This migration inserts 1000 random subscriptions for testing purposes

DO $$
DECLARE
    i INTEGER;
    random_service VARCHAR(255);
    random_price INTEGER;
    random_user_id UUID;
    random_start_date DATE;
    random_end_date DATE;
    services VARCHAR(255)[] := ARRAY['Netflix', 'Spotify', 'YouTube Premium', 'Amazon Prime', 'Disney+', 'HBO Max', 'Apple TV+', 'Hulu', 'Paramount+', 'Peacock', 'Tidal', 'Deezer', 'Yandex Music', 'VK Music', 'iCloud+', 'Google One', 'Microsoft 365', 'Adobe Creative Cloud', 'Figma', 'Notion', 'Slack', 'Discord Nitro', 'Telegram Premium', 'PlayStation Plus', 'Xbox Game Pass', 'Nintendo Switch Online', 'Steam', 'Epic Games', 'Origin', 'Uplay', 'EA Play', 'GeForce Now', 'Shadow PC', 'Parsec', 'Zoom', 'Google Meet', 'Webex', 'GoToMeeting', 'BlueJeans', 'Loom', 'Loomio', 'Calendly', 'ScheduleOnce', 'Doodle', 'When2meet', 'SurveyMonkey', 'Typeform', 'Google Forms', 'Microsoft Forms', 'JotForm', 'Zoho Forms', '123FormBuilder', 'Wufoo', 'Paperform', 'Cognito Forms', 'Formstack', 'Airtable', 'Notion', 'Coda', 'Quip', 'Evernote', 'OneNote', 'Bear', 'Ulysses', 'Scrivener', 'Grammarly', 'ProWritingAid', 'Hemingway Editor', 'Readable', 'Copyscape', 'Turnitin', 'Plagiarism Checker', 'DupliChecker', 'SmallSEOTools', 'SEMrush', 'Ahrefs', 'Moz', 'Majestic', 'Serpstat', 'SpyFu', 'KWFinder', 'LongTailPro', 'UberSuggest', 'Google Keyword Planner', 'Bing Webmaster Tools', 'Yandex Webmaster', 'Baidu Webmaster', 'Google Analytics', 'Mixpanel', 'Amplitude', 'Heap', 'Hotjar', 'Crazy Egg', 'FullStory', 'LogRocket', 'Smartlook', 'Mouseflow', 'Lucky Orange', 'Inspectlet', 'SessionCam', 'Glassbox', 'Contentsquare', 'Decibel', 'Quantum Metric', 'Pendo', 'WalkMe', 'Appcues', 'Userpilot', 'UserGuiding', 'Whatfix', 'Stonly', 'Help Scout', 'Zendesk', 'Freshdesk', 'Intercom', 'Drift', 'Crisp', 'Tidio', 'LiveChat', 'Olark', 'Tawk.to', 'Chatra', 'Formilla', 'Pure Chat', 'SnapEngage', 'Comm100', 'Kayako', 'Freshchat', 'LiveAgent', 'Zoho Desk', 'HubSpot Service', 'Salesforce Service', 'ServiceNow', 'Jira Service', 'Freshservice', 'SolarWinds', 'ManageEngine', 'NinjaRMM', 'Atera', 'Syncro', 'ConnectWise', 'Kaseya', 'N-able', 'Pulseway', 'Datto', 'Continuum', 'Addigy', 'Jamf', 'Kandji', 'Mosyle', 'Filewave', 'Fleetsmith', 'Workspace ONE', 'Intune', 'SCCM', 'Chef', 'Puppet', 'Ansible', 'SaltStack', 'Terraform', 'Pulumi', 'CloudFormation', 'Azure Resource Manager', 'Google Cloud Deployment', 'AWS CDK', 'Serverless Framework', 'SAM', 'Amplify', 'Firebase', 'Supabase', 'Appwrite', 'Nhost', 'Hasura', 'Postgraphile', 'Prisma', 'Drizzle', 'TypeORM', 'Sequelize', 'MikroORM', 'Objection.js', 'Bookshelf', 'Knex', 'Kysely'];
    user_ids UUID[] := ARRAY[
        '550e8400-e29b-41d4-a716-446655440000'::UUID,
        '550e8400-e29b-41d4-a716-446655440001'::UUID,
        '550e8400-e29b-41d4-a716-446655440002'::UUID,
        '550e8400-e29b-41d4-a716-446655440003'::UUID,
        '550e8400-e29b-41d4-a716-446655440004'::UUID,
        '550e8400-e29b-41d4-a716-446655440005'::UUID,
        '550e8400-e29b-41d4-a716-446655440006'::UUID,
        '550e8400-e29b-41d4-a716-446655440007'::UUID,
        '550e8400-e29b-41d4-a716-446655440008'::UUID,
        '550e8400-e29b-41d4-a716-446655440009'::UUID
    ];
BEGIN
    FOR i IN 1..1000 LOOP
        random_service := services[1 + floor(random() * array_length(services, 1))::int];
        random_price := 100 + floor(random() * 9900)::int; -- Price between 100 and 10000
        random_user_id := user_ids[1 + floor(random() * array_length(user_ids, 1))::int];
        random_start_date := CURRENT_DATE - (floor(random() * 365 * 2)::int || ' days')::interval; -- Last 2 years
        
        -- 30% chance to have an end date
        IF random() < 0.3 THEN
            random_end_date := random_start_date + (floor(random() * 365 + 30)::int || ' days')::interval;
        ELSE
            random_end_date := NULL;
        END IF;
        
        INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date, created_at, updated_at)
        VALUES (
            gen_random_uuid(),
            random_service,
            random_price,
            random_user_id,
            random_start_date,
            random_end_date,
            NOW(),
            NOW()
        );
    END LOOP;
END $$;
