package animal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Service struct {
	animalRepo animalRepo.RepositoryInterface
}

func NewAnimalService(
	animalRepo animalRepo.RepositoryInterface,
) *Service {
	return &Service{
		animalRepo: animalRepo,
	}
}

func (s *Service) CreateAnimal(ctx context.Context, input animal) error {

	channel, err := s.channelRepo.FindChannelByID(ctx, input.Channel, channelRepo.WithAgentPreload())
	if err != nil {
		return fmt.Errorf("failed to find channel: %w", err)
	}

	if !channel.Agent.Published {
		log.Info().
			Uint("agent_id", channel.AgentID).
			Msg("agent is not published - skipping message processing")
		return nil
	}

	phoneNumber := input.Data.Key.ExtractPhoneNumber()
	if phoneNumber == "" {
		log.Error().
			Str("remote_jid", input.Data.Key.RemoteJid).
			Msg("invalid remote JID format - cannot extract phone number")
		return nil
	}

	if channel.ConnectionStatus != "open" {
		return fmt.Errorf("channel connection status not avaiable: %s", channel.ConnectionStatus)
	}

	lead, err := s.leadRepo.FindByPhoneNumberAndWorkspace(phoneNumber, channel.Agent.WorkspaceID)
	if err != nil {
		log.Info().Msgf("lead not found for phone number: %s and workspace ID: %d", phoneNumber, channel.Agent.WorkspaceID)
		log.Info().Msgf("creating new lead for phone number: %s and workspace ID: %d", phoneNumber, channel.Agent.WorkspaceID)

		newLead := &models.Lead{
			Name:        input.Data.PushName,
			PhoneNumber: phoneNumber,
			WorkspaceID: channel.Agent.WorkspaceID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Origin:      string(ChannelWhatsApp),
		}

		err = s.leadRepo.Create(newLead)
		if err != nil {
			return fmt.Errorf("failed to create lead: %w", err)
		}
		lead = newLead
	}

	visitorID := models.ChatLeadPrefix + phoneNumber

	log.Info().
		Str("visitor_id", visitorID).
		Str("agent_hash_id", channel.Agent.HashID).
		Uint("lead_id", lead.ID).
		Uint("channel_id", channel.ID).
		Msg("creating or retrieving chat")

	chat, err := s.GetOrCreateChat(ctx, channel.Agent.HashID, visitorID, &lead.ID, &channel.ID, ChannelWhatsApp)
	if err != nil {
		return fmt.Errorf("failed to get or create chat: %w", err)
	}

	if input.Data.Key.FromMe {
		log.Info().
			Str("message_id", input.Data.Key.ID).
			Msg("message is a human intervention from WhatsApp")

		shouldUpdateChat := chat.AgentActive || chat.User == nil

		if shouldUpdateChat {
			log.Info().
				Str("session_id", chat.SessionID).
				Bool("current_agent_active", chat.AgentActive).
				Msg("human intervention detected - disabling AI agent and assigning main user")

			mainUser, err := s.userRepo.FindMainUserByWorkspaceID(channel.Agent.WorkspaceID)
			if err != nil {
				return fmt.Errorf("failed to find main user by workspace ID: %w", err)
			}

			updateFields := map[string]interface{}{
				"user_id":      mainUser.ID,
				"agent_active": false,
			}

			err = s.chatRepo.UpdateChatFields(chat.ID, updateFields)
			if err != nil {
				return fmt.Errorf("failed to update chat for WhatsApp human intervention: %w", err)
			}
		}
	}

	var sendMessageInput SendMessageInput
	switch input.Data.MessageType {
	case MessageTypeConversation:
		var conversationEventMsg MessageConversation
		if err := json.Unmarshal(input.Data.Message, &conversationEventMsg); err != nil {

			log.Error().
				Err(err).
				Str("message_raw", string(input.Data.Message)).
				Msg("failed to unmarshal conversation message")
			return nil
		}

		if conversationEventMsg.Conversation == "" {
			log.Warn().
				Str("message_id", input.Data.Key.ID).
				Msg("empty conversation message - skipping")
			return nil
		}

		log.Info().
			Str("conversation_text", conversationEventMsg.Conversation).
			Str("message_id", input.Data.Key.ID).
			Msg("processing conversation message")

		sendMessageInput = SendMessageInput{
			Message:   conversationEventMsg.Conversation,
			Key:       input.Data.Key,
			Timestamp: input.Data.Timestamp,
		}

	default:
		log.Warn().
			Str("message_type", string(input.Data.MessageType)).
			Msg("unsupported message type - skipping")
		return nil
	}

	_, err = s.SendMessage(ctx, chat.SessionID, sendMessageInput, ChannelWhatsApp)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
